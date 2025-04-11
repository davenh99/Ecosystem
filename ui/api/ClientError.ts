export class ClientError extends Error {
  url: string = "";
  status: number = 0;
  response: { [key: string]: any } = {};
  //   isAbort: boolean = false;
  originalError: any = null;

  constructor(err?: any) {
    super("ClientError");

    Object.setPrototypeOf(this, ClientError.prototype);

    // TODO figure out below and understand it... and make a new version that isnt so plagiarism
    if (err !== null && typeof err === "object") {
      this.url = typeof err.url === "string" ? err.url : "";
      this.status = typeof err.status === "number" ? err.status : 0;
      this.originalError = err.originalError;

      if (err.response !== null && typeof err.response === "object") {
        this.response = err.response;
      } else if (err.data !== null && typeof err.data === "object") {
        this.response = err.data;
      } else {
        this.response = {};
      }
    }

    this.name = "ClientError " + this.status;
    // TODO why wasn't this done in pocketbase? ( || err.message)
    this.message = this.response?.message || err.message;

    if (!this.message) {
      // TODO what could I say?
      this.message = "Error while processing request";
    }
  }
}
