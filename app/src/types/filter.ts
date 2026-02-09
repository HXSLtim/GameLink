export interface FilterOption {
  label: string
  value: string | number | boolean
}

export interface FilterSection {
  key: string
  label: string
  options: FilterOption[]
  multiple?: boolean
}

export type FilterValues = Record<
  string,
  string | number | boolean | Array<string | number | boolean>
>
