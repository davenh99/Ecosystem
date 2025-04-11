export interface CookieOpts {
  httpOnly?: boolean;
  encode?: (val: string | number | boolean) => string;
  maxAge?: number;
  domain?: string;
  path?: string;
  expires?: Date;
  secure?: boolean;
  priority?: string;
  sameSite?: boolean | "none" | "lax" | "strict";
  partitioned?: boolean;
}
