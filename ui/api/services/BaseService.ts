import Api from "../Api";

export abstract class BaseService {
  readonly api: Api;

  constructor(api: Api) {
    this.api = api;
  }
}
