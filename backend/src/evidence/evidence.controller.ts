import {
  Body,
  Controller,
  Get,
  Param,
  Post,
  Query,
} from '@nestjs/common';

import { EvidenceService } from './evidence.service';
import { IngestEvidenceDto } from './dto/ingest-evidence.dto';

@Controller('evidence')
export class EvidenceController {
  constructor(
    private readonly evidenceService: EvidenceService,
  ) {}

  @Post()
  ingest(@Body() dto: IngestEvidenceDto) {
    return this.evidenceService.ingest(dto);
  }

  @Get()
  findAll(
    @Query('endpointId') endpointId?: string,
    @Query('type') type?: string,
  ) {
    return this.evidenceService.findAll(endpointId, type);
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.evidenceService.findOne(id);
  }
}