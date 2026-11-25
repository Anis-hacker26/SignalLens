import {
  IsEnum,
  IsOptional,
  IsString,
  MaxLength,
} from 'class-validator';

export enum AssessmentClassification {
  UNDETERMINED = 'UNDETERMINED',
  BENIGN = 'BENIGN',
  SUSPICIOUS = 'SUSPICIOUS',
  MALICIOUS = 'MALICIOUS',
}

export enum AssessmentConfidence {
  LOW = 'LOW',
  MEDIUM = 'MEDIUM',
  HIGH = 'HIGH',
}

export class CreateAssessmentDto {
  @IsEnum(AssessmentClassification)
  classification!: AssessmentClassification;

  @IsEnum(AssessmentConfidence)
  confidence!: AssessmentConfidence;

  @IsOptional()
  @IsString()
  @MaxLength(10000)
  reasoning?: string;
}