import { Component, createSignal } from "solid-js";

import Header from "../layouts/header";
import TablesView from "../layouts/dashboard/TablesView";
import ListItem from "../components/listitem";

const items = ["Tables", "Fields", "Views", "Scripts", "__divider__", "Logs", "Settings"];

const Dashboard: Component = () => {
  const [screen, setScreen] = createSignal<number>(0);

  return (
    <div style={{ height: "100vh", width: "100vw" }}>
      <Header />
      <div style={{ display: "flex", "flex-direction": "row", height: "100%" }}>
        <div
          style={{
            width: "250px",
            "background-color": "var(--scheme-light)",
            color: "var(--scheme-dark)",
          }}
        >
          {items.map((item, ind) => {
            if (item == "__divider__") {
              return <div>-----</div>;
            } else {
              return (
                <ListItem onclick={() => setScreen(ind)}>
                  <p style={{ "text-decoration": ind == screen() ? "underline" : "none" }}>{item}</p>
                </ListItem>
                // <div onClick={() => setScreen(ind)}>
                //   <p style={{ "text-decoration": ind == screen() ? "underline" : "none" }}>{item}</p>
                // </div>
              );
            }
          })}
        </div>
        {/* <div>{items[screen()]}</div> */}
        <TablesView />
      </div>
    </div>
  );
};

export default Dashboard;
