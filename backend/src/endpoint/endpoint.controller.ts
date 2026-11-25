import {
  Controller,
  Get,
  Param,
} from '@nestjs/common';

import { EndpointService } from './endpoint.service';

@Controller('endpoints')
export class EndpointController {
  constructor(
    private readonly endpointService: EndpointService,
  ) {}

  @Get()
  findAll() {
    return this.endpointService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.endpointService.findOne(id);
  }
}
