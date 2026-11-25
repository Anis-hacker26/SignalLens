import { IsUUID } from 'class-validator';

export class AttachEvidenceDto {
  @IsUUID()
  evidenceId!: string;
}