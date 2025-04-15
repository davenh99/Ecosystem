import { createSignal, ParentComponent } from "solid-js";

interface Props {
  onclick: () => void;
  style?: any;
}

const ListItem: ParentComponent<Props> = (props) => {
  const [style, setStyle] = createSignal<any>(props.style);

  return (
    <div
      onclick={props.onclick}
      onmouseenter={() => setStyle({ ...props.style, "background-color": "#ccc" })}
      onmouseleave={() => setStyle({ ...props.style })}
      style={style()}
    >
      {props.children}
    </div>
  );
};

export default ListItem;
