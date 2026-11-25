import {
  BadRequestException,
  Inject,
  Injectable,
  NotFoundException,
} from '@nestjs/common';

import { DATABASE } from '../infrastructure/database/database.provider';
import { CreateInvestigationDto } from './dto/create-investigation.dto';
import { CreateAssessmentDto } from './dto/create-assessment.dto';

@Injectable()
export class InvestigationService {
  constructor(
    @Inject(DATABASE)
    private readonly db: any,
  ) {}

  async create(dto: CreateInvestigationDto) {
    const endpoint = await this.db.orm.public.Endpoint
      .where({ id: dto.endpointId })
      .first();

    if (!endpoint) {
      throw new NotFoundException(
        `Endpoint ${dto.endpointId} was not found`,
      );
    }

    return this.db.orm.public.Investigation.create({
      endpointId: dto.endpointId,
      title: dto.title,
      description: dto.description,
      status: 'OPEN',
      startedAt: new Date(dto.startedAt),
    });
  }

async close(investigationId: string) {
  const investigation =
    await this.db.orm.public.Investigation
      .where({ id: investigationId })
      .first();

  if (!investigation) {
    throw new NotFoundException(
      `Investigation ${investigationId} was not found`,
    );
  }

  if (investigation.status === 'CLOSED') {
    throw new BadRequestException(
      'Investigation is already closed',
    );
  }

  return this.db.orm.public.Investigation
    .where({ id: investigationId })
    .update({
      status: 'CLOSED',
      closedAt: new Date(),
    });
}

async findAll() {
  return this.db.orm.public.Investigation.all();
}

async findOne(id: string) {
  const investigation =
    await this.db.orm.public.Investigation
      .where({ id })
      .first();

  if (!investigation) {
    throw new NotFoundException(
      `Investigation ${id} was not found`,
    );
  }

  const endpoint =
    await this.db.orm.public.Endpoint
      .where({ id: investigation.endpointId })
      .first();

  const evidenceLinks =
    await this.db.orm.public.InvestigationEvidence
      .where({ investigationId: id })
      .all();

const evidence: any[] = [];

for (const link of evidenceLinks) {
  const item =
    await this.db.orm.public.Evidence
      .where({ id: link.evidenceId })
      .first();

  if (!item) {
    continue;
  }

  if (item.type === 'PROCESS') {
    const processEvent =
      await this.db.orm.public.ProcessEvent
        .where({ evidenceId: item.id })
        .first();

    evidence.push({
      ...item,
      processEvent,
    });

    continue;
  }

  if (item.type === 'NETWORK') {
    const networkEvent =
      await this.db.orm.public.NetworkEvent
        .where({ evidenceId: item.id })
        .first();

    evidence.push({
      ...item,
      networkEvent,
    });

    continue;
  }

  evidence.push(item);
}

  const assessment =
    await this.db.orm.public.Assessment
      .where({ investigationId: id })
      .first();

  return {
    investigation,
    endpoint,
    evidence,
    assessment,
  };
}

  async attachEvidence(
  investigationId: string,
  evidenceId: string,
) {
const investigation =
  await this.db.orm.public.Investigation
    .where({ id: investigationId })
    .first();

if (!investigation) {
  throw new NotFoundException(
    `Investigation ${investigationId} was not found`,
  );
}

if (investigation.status === 'CLOSED') {
  throw new BadRequestException(
    'Cannot attach evidence to a closed investigation',
  );
}

  const evidence =
    await this.db.orm.public.Evidence
      .where({ id: evidenceId })
      .first();

  if (!evidence) {
    throw new NotFoundException(
      `Evidence ${evidenceId} was not found`,
    );
  }

if (investigation.endpointId !== evidence.endpointId) {
  throw new BadRequestException(
    'Evidence belongs to a different endpoint than the investigation',
  );
}

const existingLink =
  await this.db.orm.public.InvestigationEvidence
    .where({
      investigationId,
      evidenceId,
    })
    .first();

if (existingLink) {
  throw new BadRequestException(
    'Evidence is already attached to this investigation',
  );
}

return this.db.orm.public.InvestigationEvidence.create({
  investigationId,
  evidenceId,
});
}

async getEvidence(investigationId: string) {
  const investigation =
    await this.db.orm.public.Investigation
      .where({ id: investigationId })
      .first();

  if (!investigation) {
    throw new NotFoundException(
      `Investigation ${investigationId} was not found`,
    );
  }

  const evidenceLinks =
    await this.db.orm.public.InvestigationEvidence
      .where({ investigationId })
      .all();

  const evidence: any[] = [];

  for (const link of evidenceLinks) {
    const item =
      await this.db.orm.public.Evidence
        .where({ id: link.evidenceId })
        .first();

    if (!item) {
      continue;
    }

    const baseEvidence = {
      ...item,
      investigationId: link.investigationId,
      addedAt: link.addedAt,
      note: link.note,
    };

    if (item.type === 'PROCESS') {
      const processEvent =
        await this.db.orm.public.ProcessEvent
          .where({ evidenceId: item.id })
          .first();

      evidence.push({
        ...baseEvidence,
        processEvent,
      });

      continue;
    }

    if (item.type === 'NETWORK') {
      const networkEvent =
        await this.db.orm.public.NetworkEvent
          .where({ evidenceId: item.id })
          .first();

      evidence.push({
        ...baseEvidence,
        networkEvent,
      });

      continue;
    }

    evidence.push(baseEvidence);
  }

  return evidence;
}

async createAssessment(
  investigationId: string,
  dto: CreateAssessmentDto,
) {
  const investigation =
    await this.db.orm.public.Investigation
      .where({ id: investigationId })
      .first();

  if (!investigation) {
    throw new NotFoundException(
      `Investigation ${investigationId} was not found`,
    );
  }

  if (investigation.status === 'CLOSED') {
  throw new BadRequestException(
    'Cannot create an assessment for a closed investigation',
  );
}

  return this.db.orm.public.Assessment.create({
    investigationId,
    classification: dto.classification,
    confidence: dto.confidence,
    summary: dto.reasoning ?? '',
  });
}

async updateAssessment(
  investigationId: string,
  dto: CreateAssessmentDto,
) {
  const investigation =
    await this.db.orm.public.Investigation
      .where({ id: investigationId })
      .first()

  if (!investigation) {
    throw new NotFoundException(
      `Investigation ${investigationId} was not found`,
    )
  }

  if (investigation.status !== 'OPEN') {
    throw new BadRequestException(
      'Closed investigations cannot be modified',
    )
  }

  const assessment =
    await this.db.orm.public.Assessment
      .where({ investigationId })
      .first()

  if (!assessment) {
    throw new NotFoundException(
      `Assessment for investigation ${investigationId} was not found`,
    )
  }

  return this.db.orm.public.Assessment
    .where({ investigationId })
    .update({
      classification: dto.classification,
      confidence: dto.confidence,
      summary: dto.reasoning ?? '',
    })
}

async getAssessment(investigationId: string) {
  const investigation =
    await this.db.orm.public.Investigation
      .where({ id: investigationId })
      .first();

  if (!investigation) {
    throw new NotFoundException(
      `Investigation ${investigationId} was not found`,
    );
  }

  const assessment =
    await this.db.orm.public.Assessment
      .where({ investigationId })
      .first();

  if (!assessment) {
    throw new NotFoundException(
      `Assessment for investigation ${investigationId} was not found`,
    );
  }

  return assessment;
}
}
