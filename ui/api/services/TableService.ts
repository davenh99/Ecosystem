import Api from "../Api";
import { TableModel } from "../utils/types";
import { CrudService } from "./CrudService";

export class TableService extends CrudService {
  constructor(api: Api) {
    super(api);
  }

  public async getList() {
    return super._getList<TableModel[]>("/tables");
  }
}
