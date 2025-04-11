import { BaseService } from "./BaseService";

export abstract class CrudService extends BaseService {
  protected async _getList<T>(urlSuffix: string): Promise<T> {
    const url: string = this.api.baseURL + urlSuffix;

    // TODO use an options object here instead of the "GET"
    return this.api.send(url, "GET") as T;
  }

  protected async _getOne<T>(urlSuffix: string, payload?: { [key: string]: any } | FormData): Promise<T> {
    const url: string = this.api.baseURL + urlSuffix;

    // TODO use an options object here instead of the "GET"
    return this.api.send(url, "GET", payload) as T;
  }

  protected async _create<T>(urlSuffix: string, payload?: { [key: string]: any } | FormData): Promise<T> {
    const url: string = this.api.baseURL + urlSuffix;

    // TODO use an options object here instead of the "POST"
    return this.api.send(url, "POST", payload) as T;
  }

  protected async _update<T>() {}

  protected async _delete<T>() {}
}
