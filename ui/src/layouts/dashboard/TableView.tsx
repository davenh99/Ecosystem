import { Component, createEffect, createSignal } from "solid-js";

import { client } from "../../../api";
import { TableModel } from "../../../api/utils/types";
import DataList from "../../components/datalist";

interface Props {
  table: TableModel;
}

const TableView: Component<Props> = (props) => {
  const [records, setRecords] = createSignal<TableModel[]>([]);

  createEffect(async () => {
    try {
      const tables = await client.table(props.table.name).getList();
      setRecords(tables);
    } catch (e) {
      alert(e);
    }
  });

  const onClickItem = (record: any) => {
    alert(record);
  };

  return <DataList model={props.table} onClickItem={onClickItem} data={records()} />;
};

export default TableView;
