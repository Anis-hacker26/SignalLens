import { Module } from '@nestjs/common';

import { DatabaseModule } from '../infrastructure/database/database.module';
import { InvestigationController } from './investigation.controller';
import { InvestigationService } from './investigation.service';

@Module({
  imports: [DatabaseModule],
  controllers: [InvestigationController],
  providers: [InvestigationService],
})
export class InvestigationModule {}