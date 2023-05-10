import {GraphRequest} from "./graph-request";
import {RequestStatus} from "./request-status";

export interface GraphBatch {
  id: Number,
  number: Number,
  status: RequestStatus,
  deleted: Date,
  graph_ids: Number[],
  base: GraphRequest
}
