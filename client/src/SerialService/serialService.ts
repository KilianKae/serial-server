import axios from 'axios';

export interface IStatus {
  name: string;
  baud: number;
  readTimeout: number;
  size: number;
  error: string;
}

const instance = axios.create({
  baseURL: 'http://localhost:8080',
  timeout: 1000,
  headers: { 'X-Custom-Header': 'foobar' },
});

export async function getStatus(): Promise<IStatus> {
  let status = (await instance.get('/api/status')).data as IStatus;
  console.log(status);
  return status;
}
