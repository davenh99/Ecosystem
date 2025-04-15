import { Component } from "solid-js";
import styles from "./header.module.css";
import { A } from "@solidjs/router";

import { user } from "../../../api";
import IconProfile from "../../components/icons/IconProfile";

const Header: Component = () => {
  return (
    <header class={styles["header-content"]}>
      <div class={styles["header-content-left"]}>
        <h4>
          <A href="/">Dave's ERP</A>
        </h4>
        <div class={styles["header-links"]}>
          <A href="/dashboard">Dashboard</A>
        </div>
        <div class={styles["header-links"]}>
          <A href="/ui">ui</A>
        </div>
      </div>
      <div>
        <IconProfile />
      </div>
    </header>
  );
};

export default Header;
