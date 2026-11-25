import {
  Inject,
  Injectable,
} from '@nestjs/common';

import { DATABASE } from '../infrastructure/database/database.provider';

@Injectable()
export class EndpointService {
  constructor(
    @Inject(DATABASE)
    private readonly db: any,
  ) {}

  async findAll() {
    return this.db.orm.public.Endpoint.all();
  }

  async findOne(id: string) {
    return this.db.orm.public.Endpoint
      .where({ id })
      .first();
  }
}
