import React from 'react'

export default function AddItemBar(props) {
  return (
    <form className="AddItemBar" onSubmit={props.onSubmit}>
      <span className="AddItemBarLogo">+</span>
      <input
        className="AddItemBarTitle"
        type="text"
        value={props.value}
        placeholder="New item"
        onChange={props.onChange}
      />
      <input className="AddItemBarButton" type="submit" value="" disabled={props.disabled} />
    </form>
  );
}