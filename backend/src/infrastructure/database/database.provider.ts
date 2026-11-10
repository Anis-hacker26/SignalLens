import { contractFromJson } from '@prisma/orm-postgres/static';
import postgres from '@prisma/orm-postgres/runtime';

import type { Contract } from '../../../prisma/contract.d.ts';
import contractJson from '../../../prisma/contract.json';

import { ConfigService } from '../config/config.service';

const contract = contractFromJson<Contract>(contractJson);

export const DATABASE = 'DATABASE';

export const databaseProvider = {
  provide: DATABASE,
  inject: [ConfigService],
  useFactory: (configService: ConfigService) => {
    const databaseUrl = configService.getDatabaseUrl();

    return postgres({
      contract,
      url: databaseUrl,
    });
  },
};