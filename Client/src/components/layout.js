import React from 'react';

import { makeStyles } from '@material-ui/core/styles';

import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';

import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';
import CircularProgress from '@material-ui/core/CircularProgress';

import PhonelinkLockIcon from '@material-ui/icons/PhonelinkLock';

import { UserContext } from "../context/user";
import { useAuthAPI } from "../utils/auth-api";

const useStyles = makeStyles((theme) => ({
    card: {
        minWidth: 275,
    },
    wrapper: {
        margin: theme.spacing(1),
        position: 'relative',
      },
    buttonProgress: {
        position: 'absolute',
        top: '50%',
        left: '50%',
        marginTop: -8,
        marginLeft: -4,
      },
      media: {
        height: 140,
      },
      qrCode: {
          width:"100%",
        },
}));

export default function Layout(props) {
    const classes = useStyles();

    const { user } = React.useContext(UserContext)

    const [isLoading, setIsLoading] = React.useState(false);
    const [mfaData, setMfaData] = React.useState({ hasMfa: false });


    const handleRegisterMfa = () => {
        setIsLoading(true)
        useAuthAPI.registerMFA(useAuthAPI.getToken())
            .then(data => {
                setIsLoading(false)
                console.log(data)
                setMfaData({data,hasMfa: true })
                console.log("MFA Regitered Successfully")

            })
            .catch(error => {
                setIsLoading(false)
                console.log("MFA Registration Falied")
            })
    }

    return (
        <div style={{ padding: 20 }}>
            <Grid container spacing={5} justify="center">
                <Grid item xs={6}  >
                    <Card className={classes.card} variant="outlined">
                        <CardContent>
                            <Typography gutterBottom variant="h5" component="h2">
                                {props.title}
                            </Typography>
                            <Typography variant="body2" component="span">
                                <pre style={{ whiteSpace: "pre-wrap" }}>{JSON.stringify(user, null, 4)}</pre>
                            </Typography>
                        </CardContent>
                        <Divider />
                        {mfaData.hasMfa && 
                        <CardContent>
                            <Typography gutterBottom variant="h5" component="h2">
                                Google Autheticator QR
                            </Typography>
                            <img className={classes.qrCode} src={`data:image/png;base64,${mfaData.data.googleAutheticator}`} alt={"QR Code"}/>
                        </CardContent>}
                        <CardActions>
                            <div className={classes.wrapper}>
                                <Button size="small" disabled={isLoading || !useAuthAPI.isAuthenticated() } startIcon={<PhonelinkLockIcon />} onClick={handleRegisterMfa}>Add MFA</Button>
                                {isLoading && <CircularProgress size={24} className={classes.buttonProgress} />}
                            </div>
                        </CardActions>
                        
                    </Card>
                </Grid>
            </Grid>
        </div>
    );
}
