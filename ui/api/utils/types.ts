export interface BaseModel {
  id: string;
}

export type FieldType =
  | "BOOLEAN"
  | "CHAR"
  | "VARCHAR"
  | "TEXT"
  | "DATETIME"
  | "TIMESTAMP"
  | "INT"
  | "JSON"
  | "FLOAT";

export interface UserViewModel extends BaseModel {
  username: string;
}

export interface UserModel extends BaseModel {
  username: string;
  email: string;
  created: Date;
  updated: Date;
}

export interface ForeignKey {
  table: string;
  column: string;
}

// once separate fields table is implemented below should extends BaseModel
export interface FieldModel {
  name: string;
  type: FieldType;
  size?: number;
  nullable?: boolean;
  primary?: boolean;
  default?: string;
  unique?: string;
  index?: boolean;
  autoIncrement?: boolean;
  foreignKey?: ForeignKey;
}

export interface TableModel extends BaseModel {
  system: boolean;
  module: string;
  name: string;
  fields: FieldModel[];
  created: Date;
  updated: Date;
}

export interface UserCreatePayload {
  email: string;
  password: string;
}
