import { createSignal, ParentComponent } from "solid-js";

interface Props {
  onclick: () => void;
  style?: any;
}

const DataListItem: ParentComponent<Props> = (props) => {
  const [style, setStyle] = createSignal<any>(props.style);

  return (
    <tr
      onclick={props.onclick}
      onmouseenter={() => setStyle({ ...props.style, "background-color": "#ccc" })}
      onmouseleave={() => setStyle({ ...props.style })}
      style={style()}
    >
      {props.children}
    </tr>
  );
};

export default DataListItem;
