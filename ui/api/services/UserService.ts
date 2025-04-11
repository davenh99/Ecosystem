import Api from "../Api";
import { UserCreatePayload, UserModel, UserViewModel } from "../utils/types";
import { CrudService } from "./CrudService";

export interface UserAuthResponse {
  User: UserModel;
  Token: string;
}

export class UserService extends CrudService {
  constructor(api: Api) {
    super(api);
  }

  // TODO allow username as well
  public async authenticate(email: string, password: string): Promise<UserModel> {
    return this._authenticate({ email, password });
  }

  public async create(payload: UserCreatePayload): Promise<UserModel> {
    const authData = await super._create<UserAuthResponse>("/user/register", payload);
    this.api.authStore.save(authData.Token, authData.User);

    return authData.User;
  }

  public async update() {}

  public async delete() {}

  public async getList() {
    return super._getList<UserViewModel[]>("/user/list");
  }

  public async authWithCookie() {
    const authData = await super._getOne<UserAuthResponse>("/user/getone");
    this.api.authStore.save(authData.Token, authData.User);

    return authData.User;
  }

  private async _authenticate(payload?: { [key: string]: any } | FormData): Promise<UserModel> {
    const authData = await this.api.send<UserAuthResponse>(
      this.api.baseURL + "/user/authenticate",
      "POST",
      payload
    );
    this.api.authStore.save(authData.Token, authData.User);

    return authData.User;
  }
}
