import React from 'react'

import TextareaAutosize from 'react-textarea-autosize';

export default React.forwardRef((props, ref) =>
    <TextareaAutosize
        {...props}
        ref={ref}
        onKeyDown={event => onKeyDown(event, props.onEnter)}
    />
);

function onKeyDown(event, f) {
    if (event.keyCode == 13) {
        f();
        event.preventDefault();
    }
}