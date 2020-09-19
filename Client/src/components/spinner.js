import React from 'react'
import CircularProgress from '@material-ui/core/CircularProgress';

export default function Spinner(props) {
    const { show } = props;
    return (
        <div style={{
            position: "absolute",
            top:"50%",
            // height: ,
            // height: height ? height: "50%",
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
        }}>
            {show ? (<CircularProgress />) : null}
        </div>
    );
}