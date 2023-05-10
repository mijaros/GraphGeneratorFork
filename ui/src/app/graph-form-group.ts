import {AbstractControl, FormControl, FormGroup, ValidationErrors, ValidatorFn, Validators} from "@angular/forms";
import {GraphType} from "./graph-type";
import {LimitsProviderService} from "./limits-provider.service";


export class GraphFormGroup {
  private node_degree: FormControl<Number> = new FormControl<Number>({value: 0, disabled: true}, {nonNullable: true})
  private node_degree_max: FormControl<Number> = new FormControl<Number>({
    value: 0,
    disabled: true
  }, {nonNullable: true})
  private node_degree_average: FormControl<Number> = new FormControl<Number>({
    value: 0,
    disabled: true
  }, {nonNullable: true})
  private connected: FormControl<Boolean> = new FormControl<boolean>(true, {nonNullable: true})

  constructor(private limits: LimitsProviderService) {
    this.formGroup.controls['weighted'].valueChanges.subscribe(v => {
      if (v) {
        this.formGroup.get('weight_min')?.enable()
        this.formGroup.get('weight_max')?.enable()
      } else {
        this.formGroup.get('weight_min')?.disable()
        this.formGroup.get('weight_max')?.disable()
      }
    })
    this.formGroup.controls['type'].valueChanges.subscribe(value1 => {
      this.node_degree_average.disable()
      this.node_degree_max.disable()
      this.node_degree.disable()
      this.connected.enable()
      switch (value1) {
        case GraphType.K_REGULAR:
        case GraphType.AT_LEAST:
          this.node_degree.enable()
          break
        case GraphType.BETWEEN:
          this.node_degree.enable()
          this.node_degree_max.enable()
          break
        case GraphType.AVERAGE:
          this.node_degree_average.enable()
          break
        case GraphType.COMPLETE:
          this.connected.disable()
      }
    })
  }

  get value() {
    return this.formGroup.value
  }

  get form() {
    return this.formGroup
  }

  get valid() {
    return this.formGroup.valid
  }

  GraphValidator: ValidatorFn = (control: AbstractControl): ValidationErrors | null => {
    let res: ValidationErrors = {}
    let type = control.get('type')?.value
    if (type == null) {
      res = {typeInvalid: true, typeInvalidMessage: "Type mustn't be null."}
    }
    let nodes = control.get('nodes')?.value
    if (nodes == null) {
      res = {...res, nodesMessage: "Number of nodes mustn't be empty.", nodesError: true}
    } else if (nodes <= 0) {
      res = {...res, nodesMessage: "Number of nodes must be greater than 0.", nodesError: true}
    } else if (nodes > this.limits.maxNodes) {
      res = {
        ...res,
        nodesMessage: `Number of nodes must be lower or equal than ${this.limits.maxNodes}.`,
        nodesError: true
      }
    } else if (!this.checkWholeNumber(nodes)) {
      res = {...res, nodesMessage: "Number of nodes must be whole number.", nodesError: true}
    }
    switch (control.get("type")?.value) {
      case GraphType.K_REGULAR:
        res = {...res, ...this.validateKRegularGraph(control)}
        break
      case GraphType.BETWEEN:
        res = {...res, ...this.validateBetweenRegularGraph(control)}
        break
      case GraphType.COMPLETE:
        break
      case GraphType.AT_LEAST:
        res = {...res, ...this.validateAtLeastGraph(control)}
        break
      case GraphType.AVERAGE:
        res = {...res, ...this.validateAverageGraph(control)}
        break
      default:
        res = {...res, invalidGraphType: true}
    }

    let weighted = control.get('weighted')?.value
    if (weighted) {
      res = {...res, ...this.validateWeightedGraph(control)}
    }
    return res
  }

  private formGroup = new FormGroup({
    type: new FormControl<GraphType | null>(null, {nonNullable: true, validators: Validators.required}),
    nodes: new FormControl<Number>(1, {nonNullable: true, validators: Validators.required}),
    node_degree: this.node_degree,
    node_degree_max: this.node_degree_max,
    node_degree_average: this.node_degree_average,
    weighted: new FormControl<boolean>(false, {nonNullable: true}),
    connected: this.connected,
    weight_min: new FormControl<Number>({value: 0, disabled: true}, {nonNullable: true}),
    weight_max: new FormControl<Number>({value: 0, disabled: true}, {nonNullable: true}),
    seed: new FormControl<Number>({value: 0, disabled: true}, {nonNullable: true})
  }, {validators: this.GraphValidator})

  reset() {
    this.formGroup.reset()
    this.formGroup.get('seed')?.disable()
  }

  checkWholeNumber(n: Number): boolean {
    if (n == null) {
      return false
    }
    return n === Math.floor(n.valueOf())
  }

  validateKRegularGraph(control: AbstractControl): ValidationErrors | null {
    let res = {}
    let nodes = control.get('nodes')?.value
    let degree = control.get('node_degree')?.value
    let connected = control.get('connected')?.value
    if (degree <= 0) {
      res = {degreeMessage: "Degree of graph must be higher than 0.", invalidDegree: true}
    } else if (degree >= nodes) {
      res = {degreeMessage: "Degree of graph must be lower than number of nodes.", invalidDegree: true}
    } else if (((nodes * degree) % 2) == 1) {
      res = {...res, degreeMessage: "Degree must be even for odd number of nodes.", invalidDegree: true}
    } else if (connected && degree < 2 && nodes > 2) {
      res = {...res, degreeMessage: "Degree must be at least 2 for connected graph.", invalidDegree: true}
    }

    return res
  }

  validateAverageGraph(control: AbstractControl): ValidationErrors | null {
    let res = {}
    let nodes = control.get('nodes')?.value as number
    let avgDeg = control.get('node_degree_average')?.value as number
    let connected = control.get('connected')?.value as number

    if (avgDeg < 0) {
      res = {avgDeg: true, avgDegMessage: "Average degree must be non-negative."}
    } else if (avgDeg > (nodes - 1)) {
      res = {avgDeg: true, avgDegMessage: "Average degree must be lower than number of nodes."}
    } else if (connected && avgDeg < 2 && nodes > 2) {
      res = {avgDeg: true, avgDegMessage: "Average degree must be at least two for connected graph."}
    }
    return res
  }

  validateAtLeastGraph(control: AbstractControl): ValidationErrors | null {
    let res = {}
    let nodes = control.get('nodes')?.value
    let degree = control.get('node_degree')?.value

    if (degree < 0) {
      res = {invalidDegree: true, degreeMessage: "Minimal degree must be non-negative number"}
    } else if (degree >= nodes) {
      res = {invalidDegree: true, degreeMessage: "Minimal degree must be lower than number of nodes"}
    }

    return res
  }

  validateBetweenRegularGraph(control: AbstractControl): ValidationErrors | null {
    let res = {}
    let nodes = control.get('nodes')?.value as Number
    let degree_min = control.get('node_degree')?.value as Number
    let degree_max = control.get('node_degree_max')?.value as Number

    if (!this.checkWholeNumber(degree_max)) {
      res = {...res, degMax: true, degMaxMessage: "Degree must be whole number."}
    }

    if (!this.checkWholeNumber(degree_min)) {
      res = {...res, degMin: true, degMinMessage: "Degree must be whole number."}
    }

    if (degree_min >= degree_max) {
      res = {
        ...res, degMaxMessage: "Max degree must be strictly higher than min",
        degMinMessage: "Max degree must be strictly higher than min",
        degMin: true, degMax: true
      }
    }

    if (degree_min >= nodes) {
      res = {degMinMessage: "Minimal degree must be lower than number of nodes.", degMin: true}
    } else if (degree_min < 0) {
      res = {degMinMessage: "Minimal degree must be non-negative number.", degMin: true}
    }

    if (degree_max >= nodes) {
      res = {...res, degMaxMessage: "Maximal degree must be lower than number of nodes.", degMax: true}
    } else if (degree_max <= 0) {
      res = {...res, degMaxMessage: "Maximal degree must be higher than 0.", degMax: true}
    }

    return res
  }

  validateWeightedGraph(control: AbstractControl): ValidationErrors | null {
    let res = {}
    let weightMin = control.get('weight_min')?.value
    let weightMax = control.get('weight_max')?.value

    if (!this.checkWholeNumber(weightMax)) {
      res = {weightMax: true, weightMaxMessage: "Weight must be a whole number."}
    }

    if (!this.checkWholeNumber(weightMin)) {
      res = {...res, weightMin: true, weightMinMessage: "Weight must be a whole number"}
    }

    if (weightMax < weightMin) {
      res = {...res, weightMax: true, weightMaxMessage: "Maximal weight must be at least equal to the weight min"}
    }

    if (weightMin == 0 && weightMax == 0) {
      res = {...res, weightMax: true, weightMin: true, weightMaxMessage: ""}
    }

    return res
  }
}
