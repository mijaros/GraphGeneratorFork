import {GraphType} from "./graph-type";
import {RequestStatus} from "./request-status";

export interface GraphRequest {
  id: Number;
  type: GraphType;
  nodes: Number;
  node_degree: Number;
  node_degree_max: Number;
  connected: Boolean;
  node_degree_average: Number;
  status: RequestStatus;
  weighted: Boolean;
  weight_min: Number;
  weight_max: Number;
  deleted: Date;
  seed: Number;
}
