import { createEffect, createSignal } from "solid-js";
import { TableModel } from "../../../api/utils/types";
import DataListItem from "./DataListItem";

// import { client } from "../../../api";
// import { TableModel } from "../../../api/utils/types";

interface Props {
  model: TableModel;
  data: any[];
  onClickItem: (record: any) => void;
}

function DataList(props: Props) {
  return (
    <table style={{ "margin-top": "10px", width: "100%" }}>
      <thead>
        <tr>
          {props.model.fields.map((field) => (
            <th style={{ "text-align": "left" }}>{field.name}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {props.data.map((record) => (
          <DataListItem onclick={() => props.onClickItem(record)}>
            {props.model.fields.map((field) => (
              <td>{String(record[field.name]) || "N/A"}</td>
            ))}
          </DataListItem>
        ))}
      </tbody>
    </table>
  );
}

export default DataList;
