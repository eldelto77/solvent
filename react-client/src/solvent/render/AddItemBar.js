import React from 'react'

import { ReactComponent as Plus } from '../../icons/plus.svg'
import { ReactComponent as PlusCircle } from '../../icons/plus-circle.svg'

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
      <button className="AddItemBarButton" type="submit" value="" disabled={props.disabled}>
        {props.disabled ? <Plus /> : <PlusCircle />}
      </ button>
    </form>
  );
}