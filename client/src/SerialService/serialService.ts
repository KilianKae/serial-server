import axios from 'axios';

export interface IStatus {
  name: string;
  baud: number;
  readTimeout: number;
  size: number;
  error: string;
}

export interface IPorts {
  Name: string;
  IsUSB: string;
  VID: string;
  PID: string;
  SerialNumber: string;
  Product: string;
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

export async function getPorts(): Promise<IPorts[]> {
  let ports = (await instance.get('/api/ports')).data as IPorts[];
  console.log(ports);
  return ports;
}