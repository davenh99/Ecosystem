import { Component } from "solid-js";

import Header from "../layouts/header";

const Home: Component = () => {
  return (
    <div>
      <Header />
      <div>
        <p>You have no apps installed!</p>
      </div>
    </div>
  );
};

export default Home;
