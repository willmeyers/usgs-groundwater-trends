export interface Mouse {
  x: number;
  y: number;
}

export interface Site {
  ID: number;
  SiteName: string;
  SiteNo: string;
  Location: { type: string; coordinates: [number, number] };
  TSEstimateInMM: number;
  MKZ: number;
}

export interface Datapoint {
  Value: number;
  Year: number;
}

export interface SiteCache {
  [key: string]: Datapoint[];
}
