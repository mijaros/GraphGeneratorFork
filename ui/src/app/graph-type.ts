export enum GraphType {
  K_REGULAR = "exact-degree",
  AVERAGE = "average-degree",
  BETWEEN = "between-degree",
  AT_LEAST = "at-least-degree",
  COMPLETE = "complete"
}

export const keys = Object.keys(GraphType)
export const values: GraphType[] = [
  GraphType.K_REGULAR,
  GraphType.AVERAGE,
  GraphType.BETWEEN,
  GraphType.AT_LEAST,
  GraphType.COMPLETE]

export function explanation(t: GraphType): string {
  switch (t) {
    case GraphType.AVERAGE:
      return "Average degree"
    case GraphType.BETWEEN:
      return "Between Degree"
    case GraphType.COMPLETE:
      return "Complete"
    case GraphType.K_REGULAR:
      return "Regular"
    case GraphType.AT_LEAST:
      return "At least"
  }
}
