import React from 'react'

import TextareaAutosize from 'react-textarea-autosize';

import { ReactComponent as Plus } from '../../icons/plus.svg'
import { ReactComponent as PlusCircle } from '../../icons/plus-circle.svg'

export default function AddItemBar(props) {

  const inputField = React.createRef();

  return (
    <form className="AddItemBar" onSubmit={props.onSubmit} onClick={() => inputField.current.focus()}>
      <button className="AddItemBarButton" type="submit" value="" disabled={props.disabled}>
        {props.disabled ? <Plus /> : <PlusCircle />}
      </ button>
      <TextareaAutosize
        ref={inputField}
        className="AddItemBarTitle"
        type="text"
        placeholder="New item"
        onChange={props.onChange}
      >
        {props.value}
      </TextareaAutosize>
      <button className="AddItemBarLogo" type="button" disabled="true">+</button>
    </form>
  );
}