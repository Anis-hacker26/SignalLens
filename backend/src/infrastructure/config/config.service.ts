import { Injectable } from '@nestjs/common';
import { ConfigService as NestConfigService } from '@nestjs/config';

@Injectable()
export class ConfigService {
  constructor(private readonly config: NestConfigService) {}

  getDatabaseUrl(): string {
    const databaseUrl = this.config.get<string>('DATABASE_URL');

    if (!databaseUrl) {
      throw new Error('DATABASE_URL is not configured');
    }

    return databaseUrl;
  }

  getPort(): number {
    return this.config.get<number>('PORT', 3000);
  }

  getNodeEnv(): string {
    return this.config.get<string>('NODE_ENV', 'development');
  }
  isDatabaseConfigured(): boolean {
  return Boolean(this.config.get<string>('DATABASE_URL'));
}
getFrontendUrl(): string {
  return this.config.get<string>(
    'FRONTEND_URL',
    'http://localhost:5173',
  );
}

}