export enum RequestStatus {
  FINISHED = "finished",
  IN_PROGRESS = "not-finished",
  DELETED = "deleted",
  UNDEFINED = "undefined"
}

export function explain(g: RequestStatus): string {
  switch (g) {
    case RequestStatus.FINISHED:
      return "Finished"
    case RequestStatus.IN_PROGRESS:
      return "In progress"
    case RequestStatus.DELETED:
      return "Deleted"
    case RequestStatus.UNDEFINED:
      return "undefined"
  }
  return ""
}

export function getClass(g: RequestStatus): string {
  switch (g) {
    case RequestStatus.FINISHED:
      return "bg-success"
    case RequestStatus.IN_PROGRESS:
      return "bg-primary"
    case RequestStatus.DELETED:
      return "bg-secondary"
    case RequestStatus.UNDEFINED:
      return "bg-warn"
  }
  return ""
}
