import { Module, OnModuleDestroy, OnModuleInit, Inject } from '@nestjs/common';

import { ConfigModule } from '../config/config.module';
import { DATABASE, databaseProvider } from './database.provider';

@Module({
  imports: [ConfigModule],
  providers: [databaseProvider],
  exports: [DATABASE],
})
export class DatabaseModule implements OnModuleInit, OnModuleDestroy {
  constructor(
    @Inject(DATABASE)
    private readonly database: {
      connect: () => Promise<void>;
      close: () => Promise<void>;
    },
  ) {}

  async onModuleInit(): Promise<void> {
    await this.database.connect();
  }

  async onModuleDestroy(): Promise<void> {
    await this.database.close();
  }
}