import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';

import { ConfigModule } from './infrastructure/config/config.module';
import { DatabaseModule } from './infrastructure/database/database.module';
import { EvidenceModule } from './evidence/evidence.module';
import { InvestigationModule } from './investigation/investigation.module';
import { EndpointModule } from './endpoint/endpoint.module';

@Module({
  imports: [
  ConfigModule,
  DatabaseModule,
  EvidenceModule,
  InvestigationModule,
  EndpointModule,
],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}