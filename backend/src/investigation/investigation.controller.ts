import {
  Body,
  Controller,
  Get,
  Param,
  Patch,
  Post,
} from '@nestjs/common';

import { AttachEvidenceDto } from './dto/attach-evidence.dto';
import { CreateAssessmentDto } from './dto/create-assessment.dto';
import { CreateInvestigationDto } from './dto/create-investigation.dto';
import { InvestigationService } from './investigation.service';

@Controller('investigations')
export class InvestigationController {
  constructor(
    private readonly investigationService: InvestigationService,
  ) {}

  @Post()
  create(@Body() dto: CreateInvestigationDto) {
    return this.investigationService.create(dto);
  }

  @Get()
  findAll() {
    return this.investigationService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.investigationService.findOne(id);
  }

  @Post(':investigationId/evidence')
  attachEvidence(
    @Param('investigationId') investigationId: string,
    @Body() dto: AttachEvidenceDto,
  ) {
    return this.investigationService.attachEvidence(
      investigationId,
      dto.evidenceId,
    );
  }

  @Post(':investigationId/assessments')
createAssessment(
  @Param('investigationId') investigationId: string,
  @Body() dto: CreateAssessmentDto,
) {
  return this.investigationService.createAssessment(
    investigationId,
    dto,
  );
}

  @Get(':investigationId/evidence')
  getEvidence(
    @Param('investigationId') investigationId: string,
  ) {
    return this.investigationService.getEvidence(investigationId);
  }

  @Get(':investigationId/assessment')
getAssessment(
  @Param('investigationId') investigationId: string,
) {
  return this.investigationService.getAssessment(
    investigationId,
  );
}

@Patch(':investigationId/assessment')
updateAssessment(
  @Param('investigationId') investigationId: string,
  @Body() dto: CreateAssessmentDto,
) {
  return this.investigationService.updateAssessment(
    investigationId,
    dto,
  )
}
@Patch(':investigationId/close')
close(
  @Param('investigationId') investigationId: string,
) {
  return this.investigationService.close(investigationId);
}

}