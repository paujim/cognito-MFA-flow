import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import Container from '@material-ui/core/Container';

import { AuthContext } from "./context/auth";
import Spinner from './components/spinner';
import LoginForm from './components/loginForm'

const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
    },
    menuButton: {
        marginRight: theme.spacing(2),
    },
    title: {
        flexGrow: 1,
    },
}));



export default function ButtonAppBar() {

    const [open, setOpen] = React.useState(false);

    const handleClickOpen = () => {
        setOpen(true);
    };
    const handleClose = () => {
        setOpen(false);
    };


    const [user] = React.useState();
    const classes = useStyles();
    return (
        <AuthContext.Provider value={user}>
            <Container maxWidth="lg">
                <div className={classes.root}>
                    <AppBar position="static">
                        <Toolbar>
                            <Typography variant="h6" className={classes.title}>
                                MFA Sample
                            </Typography>
                            <Button color="inherit" onClick={handleClickOpen}>Login</Button>
                            <LoginForm open={open} handleClose={handleClose}/>
                        </Toolbar>
                    </AppBar>
                    <Spinner show={false}/>
                </div>
            </Container>
        </AuthContext.Provider>
    );
}
