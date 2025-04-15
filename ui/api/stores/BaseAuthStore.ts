import { UserModel } from "../utils/types";
// import { CookieOpts } from "../utils/cookie";

export type AuthModel = UserModel | null;
export type OnStoreChangeFunc = (token: string, user: AuthModel) => void;
const defaultCookieKey = "dave_erp_db_auth";

export class BaseAuthStore {
  protected baseToken: string = "";
  protected baseModel: AuthModel = null;
  private _onChangeCallback: OnStoreChangeFunc | null = null;

  get token(): string {
    return this.baseToken;
  }

  get user(): AuthModel {
    return this.baseModel;
  }

  public save(token: string, record: AuthModel): void {
    this.baseToken = token || "";
    this.baseModel = record || null;
    this.triggerChange();
  }

  public clear(): void {
    this.baseToken = "";
    this.baseModel = null;
    this.triggerChange();
  }

  public loadFromCookie(key = defaultCookieKey) {}

  public saveToCookie(key = defaultCookieKey) {
    // const cookieOpts: CookieOpts = {
    //   secure: true, // https
    //   httpOnly: true,
    //   sameSite: "strict",
    //   path: "/",
    //   maxAge: 31536000000, // 1 year
    // };
  }

  public deleteCookie(key = defaultCookieKey) {}

  // exportToCookie(options?: SerializeOptions, key = defaultCookieKey): string {
  //   const defaultOptions: SerializeOptions = {
  //     secure: true,
  //     sameSite: true,
  //     httpOnly: true,
  //     path: "/",
  //   };

  //   // extract the token expiration date
  //   const payload = getTokenPayload(this.token);
  //   if (payload?.exp) {
  //     defaultOptions.expires = new Date(payload.exp * 1000);
  //   } else {
  //     defaultOptions.expires = new Date("1970-01-01");
  //   }

  //   // merge with the user defined options
  //   options = Object.assign({}, defaultOptions, options);

  //   const rawData = {
  //     token: this.token,
  //     record: this.record ? JSON.parse(JSON.stringify(this.record)) : null,
  //   };

  //   let result = cookieSerialize(key, JSON.stringify(rawData), options);

  //   const resultLength = typeof Blob !== "undefined" ? new Blob([result]).size : result.length;

  //   // strip down the model data to the bare minimum
  //   if (rawData.record && resultLength > 4096) {
  //     rawData.record = { id: rawData.record?.id, email: rawData.record?.email };
  //     const extraProps = ["collectionId", "collectionName", "verified"];
  //     for (const prop in this.record) {
  //       if (extraProps.includes(prop)) {
  //         rawData.record[prop] = this.record[prop];
  //       }
  //     }
  //     result = cookieSerialize(key, JSON.stringify(rawData), options);
  //   }

  //   return result;
  // }

  public onChange(callback: OnStoreChangeFunc): void {
    this._onChangeCallback = callback;
  }

  protected triggerChange(): void {
    if (this._onChangeCallback) {
      this._onChangeCallback(this.token, this.user);
    }
  }
}
