import {
  IsBoolean,
  IsDateString,
  IsEnum,
  IsInt,
  IsObject,
  IsOptional,
  IsString,
  IsUUID,
  Min,
} from 'class-validator';

enum EvidenceType {
  PROCESS = 'PROCESS',
  NETWORK = 'NETWORK',
  DNS = 'DNS',
  FILE = 'FILE',
  POWERSHELL = 'POWERSHELL',
}

export class IngestEvidenceDto {
  @IsString()
  agentId!: string;

  @IsString()
  hostname!: string;

  @IsString()
  osName!: string;

  @IsString()
  osVersion!: string;

  @IsString()
  agentVersion!: string;

  @IsEnum(EvidenceType)
  type!: EvidenceType;

  @IsString()
  source!: string;

  @IsOptional()
  @IsString()
  sourceEventId?: string;

  @IsDateString()
  observedAt!: string;

  @IsDateString()
  collectedAt!: string;

  @IsOptional()
  @IsString()
  correlationId?: string;

  @IsObject()
  payload!: Record<string, unknown>;

  // Process fields
  @IsOptional()
  @IsString()
  processGuid?: string;

  @IsOptional()
  @IsInt()
  @Min(0)
  processId?: number;

  @IsOptional()
  @IsInt()
  @Min(0)
  parentProcessId?: number;

  @IsOptional()
  @IsString()
  parentProcessGuid?: string;

  @IsOptional()
  @IsString()
  imagePath?: string;

  @IsOptional()
  @IsString()
  commandLine?: string;

  @IsOptional()
  @IsString()
  currentDirectory?: string;

  @IsOptional()
  @IsString()
  user?: string;

  @IsOptional()
  @IsString()
  integrityLevel?: string;

  @IsOptional()
  @IsInt()
  @Min(0)
  terminalSessionId?: number;

  @IsOptional()
  @IsString()
  fileVersion?: string;

  @IsOptional()
  @IsString()
  description?: string;

  @IsOptional()
  @IsString()
  product?: string;

  @IsOptional()
  @IsString()
  company?: string;

  @IsOptional()
  @IsString()
  originalFileName?: string;

  // Network fields
  @IsOptional()
  @IsString()
  protocol?: string;

  @IsOptional()
  @IsBoolean()
  initiated?: boolean;

  @IsOptional()
  @IsString()
  sourceIp?: string;

  @IsOptional()
  @IsInt()
  @Min(0)
  sourcePort?: number;

  @IsOptional()
  @IsString()
  destinationIp?: string;

  @IsOptional()
  @IsInt()
  @Min(0)
  destinationPort?: number;

  @IsOptional()
  @IsString()
  sourceHostname?: string;

  @IsOptional()
  @IsString()
  destinationHostname?: string;

  @IsOptional()
  @IsBoolean()
  sourceIsIpv6?: boolean;

  @IsOptional()
  @IsBoolean()
  destinationIsIpv6?: boolean;
}