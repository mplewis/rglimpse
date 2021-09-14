export type Torrent = {
  completed: boolean;
  completed_bytes: number;
  created: Date;
  down_rate: number;
  finished: Date;
  hash: string;
  label: string;
  name: string;
  path: string;
  ratio: number;
  size: number;
  started: Date;
  up_rate: number;
};
