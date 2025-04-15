import Api from "../Api";
import { TableModel } from "../utils/types";
import { CrudService } from "./CrudService";

export class RecordService<M = TableModel> extends CrudService {
  table: string;

  constructor(api: Api, tableName: string) {
    super(api);

    this.table = tableName;
  }

  public async getList<T = M>() {
    return super._getList<T[]>(`/tables/${this.table}/records/`);
  }
}
