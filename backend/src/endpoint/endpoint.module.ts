import { Module } from '@nestjs/common';

import { DatabaseModule } from '../infrastructure/database/database.module';
import { EndpointController } from './endpoint.controller';
import { EndpointService } from './endpoint.service';

@Module({
  imports: [DatabaseModule],
  controllers: [EndpointController],
  providers: [EndpointService],
})
export class EndpointModule {}
