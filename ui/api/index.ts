import Api from "./Api";

export { ClientError } from "./ClientError";

export const client = new Api("http://127.0.0.1:8080/api/v1");

export let user = client.authStore.user;

client.authStore.onChange(() => {
  user = client.authStore.user;
});
