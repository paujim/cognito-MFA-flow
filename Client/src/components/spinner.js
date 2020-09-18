import React from 'react'
import CircularProgress from '@material-ui/core/CircularProgress';

export default function Spinner(props) {
    const { children, classes, show, ...other } = props;
    return (
        <div style={{
            position: "absolute",
            // height: "50%",
            // height: '90vh',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
        }}>
            {show ? (<CircularProgress />) : null}
        </div>
    );
}