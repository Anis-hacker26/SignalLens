import {
  Inject,
  Injectable,
  NotFoundException,
} from '@nestjs/common';

import { DATABASE } from '../infrastructure/database/database.provider';
import { IngestEvidenceDto } from './dto/ingest-evidence.dto';

@Injectable()
export class EvidenceService {
  constructor(
    @Inject(DATABASE)
    private readonly db: any,
  ) {}

  async ingest(dto: IngestEvidenceDto) {
    const receivedAt = new Date();
    const observedAt = new Date(dto.observedAt);
    const collectedAt = new Date(dto.collectedAt);

    const endpoint = await this.db.orm.public.Endpoint.upsert({
      conflictOn: { agentId: dto.agentId },
      create: {
        agentId: dto.agentId,
        hostname: dto.hostname,
        osName: dto.osName,
        osVersion: dto.osVersion,
        agentVersion: dto.agentVersion,
        firstSeenAt: observedAt,
        lastSeenAt: observedAt,
      },
      update: {
        hostname: dto.hostname,
        osName: dto.osName,
        osVersion: dto.osVersion,
        agentVersion: dto.agentVersion,
        lastSeenAt: observedAt,
      },
    });

    const rawTelemetry =
      await this.db.orm.public.RawTelemetry.create({
        endpointId: endpoint.id,
        source: dto.source,
        sourceEventId: dto.sourceEventId,
        observedAt,
        collectedAt,
        receivedAt,
        payload: dto.payload,
      });

    const evidence =
      await this.db.orm.public.Evidence.create({
        endpointId: endpoint.id,
        type: dto.type,
        source: dto.source,
        sourceEventId: dto.sourceEventId,
        observedAt,
        collectedAt,
        receivedAt,
        correlationId: dto.correlationId,
        rawTelemetryId: rawTelemetry.id,
      });

    if (dto.type === 'PROCESS') {
      await this.db.orm.public.ProcessEvent.create({
        evidenceId: evidence.id,
        processGuid: dto.processGuid,
        processId: dto.processId,
        parentProcessId: dto.parentProcessId,
        parentProcessGuid: dto.parentProcessGuid,
        imagePath: dto.imagePath,
        commandLine: dto.commandLine,
        currentDirectory: dto.currentDirectory,
        user: dto.user,
        integrityLevel: dto.integrityLevel,
        terminalSessionId: dto.terminalSessionId,
        fileVersion: dto.fileVersion,
        description: dto.description,
        product: dto.product,
        company: dto.company,
        originalFileName: dto.originalFileName,
      });
    }

    if (dto.type === 'NETWORK') {
      await this.db.orm.public.NetworkEvent.create({
        evidenceId: evidence.id,
        processGuid: dto.processGuid,
        processId: dto.processId,
        protocol: dto.protocol,
        initiated: dto.initiated,
        sourceIp: dto.sourceIp,
        sourcePort: dto.sourcePort,
        destinationIp: dto.destinationIp,
        destinationPort: dto.destinationPort,
        sourceHostname: dto.sourceHostname,
        destinationHostname: dto.destinationHostname,
        sourceIsIpv6: dto.sourceIsIpv6,
        destinationIsIpv6: dto.destinationIsIpv6,
      });
    }

    return {
      evidenceId: evidence.id,
      rawTelemetryId: rawTelemetry.id,
      endpointId: endpoint.id,
      receivedAt,
    };
  }

async findAll(
  endpointId?: string,
  type?: string,
) {
  const where: Record<string, unknown> = {};

  if (endpointId) {
    where.endpointId = endpointId;
  }

  if (type) {
    where.type = type;
  }

  return this.db.orm.public.Evidence
    .where(where)
    .all();
}

async findOne(id: string) {
  const evidence = await this.db.orm.public.Evidence
    .where({ id })
    .first();

  if (!evidence) {
    throw new NotFoundException(
      `Evidence ${id} was not found`,
    );
  }

  const processEvent =
    await this.db.orm.public.ProcessEvent
      .where({ evidenceId: id })
      .first();

  const networkEvent =
    await this.db.orm.public.NetworkEvent
      .where({ evidenceId: id })
      .first();

  const dnsEvent =
    await this.db.orm.public.DnsEvent
      .where({ evidenceId: id })
      .first();

  const fileEvent =
    await this.db.orm.public.FileEvent
      .where({ evidenceId: id })
      .first();

  const powerShellEvent =
    await this.db.orm.public.PowerShellEvent
      .where({ evidenceId: id })
      .first();

  return {
    ...evidence,
    processEvent,
    networkEvent,
    dnsEvent,
    fileEvent,
    powerShellEvent,
  };
}
}
