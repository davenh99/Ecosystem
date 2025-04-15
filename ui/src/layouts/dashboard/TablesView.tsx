import { createEffect, createSignal } from "solid-js";

import { client } from "../../../api";
import { TableModel } from "../../../api/utils/types";
import DataList from "../../components/datalist";
import TableView from "./TableView";

const tableTableModel: TableModel = {
  id: "",
  system: true,
  module: "",
  name: "",
  fields: [
    { name: "name", type: "TEXT" },
    { name: "system", type: "BOOLEAN" },
    { name: "module", type: "TEXT" },
    { name: "created", type: "DATETIME" },
    { name: "updated", type: "DATETIME" },
  ],
  created: new Date(),
  updated: new Date(),
};

function TablesView() {
  const [breadcrumb, setBreadcrumb] = createSignal<string>("Tables");
  const [itemClicked, setItemClicked] = createSignal<boolean>(false);
  const [clickedTable, setClickedTable] = createSignal<TableModel>(tableTableModel);
  const [tables, setTables] = createSignal<TableModel[]>([]);

  createEffect(async () => {
    try {
      const tables = await client.tables.getList();
      setTables(tables);
    } catch (e) {
      alert(e);
    }
  });

  const onClickItem = (table: TableModel) => {
    setBreadcrumb(`Tables > ${table.name}`);
    setItemClicked(true);
    setClickedTable(table);
  };

  return (
    <div style={{ padding: "20px", width: "100%" }}>
      <h2>{breadcrumb()}</h2>
      {!itemClicked() || clickedTable().name == "" ? (
        <DataList model={tableTableModel} onClickItem={onClickItem} data={tables()} />
      ) : (
        <TableView table={clickedTable()} />
      )}
    </div>
  );
}

export default TablesView;
