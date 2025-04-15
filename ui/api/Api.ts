import { BaseAuthStore } from "./stores/BaseAuthStore";
import { UserService } from "./services/UserService";
import { RecordService } from "./services/RecordService";
import { TableService } from "./services/TableService";
import { ClientError } from "./ClientError";
import { TableModel } from "./utils/types";

export default class Api {
  baseURL: string;
  authStore: BaseAuthStore;
  users: UserService;
  tables: TableService;
  private recordServices: { [key: string]: RecordService } = {};

  constructor(baseURL: string) {
    this.baseURL = baseURL;
    this.authStore = new BaseAuthStore();
    this.users = new UserService(this);
    this.tables = new TableService(this);
  }

  table<M = TableModel>(name: string): RecordService<M> {
    console.log(name);
    if (!this.recordServices[name]) {
      this.recordServices[name] = new RecordService(this, name);
    }

    return this.recordServices[name];
  }

  // redo below maybe to avoid plagiarism!!! and rename and stuff
  public async send<T>(
    url: string,
    method?: string,
    payload?: { [key: string]: any } | FormData
  ): Promise<T> {
    // TODO does below need to consider not having method or payload??
    // TODO extract out below!!! to an options object thingy, so not every request has same options
    // also, send authorization header
    const options: RequestInit = {
      method: method,
      body: JSON.stringify(payload),
      credentials: "include",
    };

    try {
      const response: Response = await fetch(url, options);
      const data = await response.json();

      if (response.status >= 400) {
        throw new ClientError({
          url: url,
          status: response.status,
          data: data,
        });
      }

      return data as T;
    } catch (e: any) {
      throw new ClientError(e);
    }
  }
}
