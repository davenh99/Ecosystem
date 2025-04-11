import { Component } from "solid-js";

import styles from "./icons.module.css";

const IconProfile: Component = () => {
  return (
    <div class={styles.profile}>
      <p class={styles["profile-text"]}>DH</p>
    </div>
  );
};

export default IconProfile;
