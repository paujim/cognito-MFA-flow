import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import { UserContext } from "../context/user";


const useStyles = makeStyles((theme) => ({
    card: {
        minWidth: 275,
    },
    bullet: {
        display: 'inline-block',
        margin: '0 2px',
        transform: 'scale(0.8)',
    },
    title: {
        fontSize: 14,
    },
    pos: {
        marginBottom: 12,
    },
}));

export default function Dashboard(props) {
    const classes = useStyles();

    const {user} = React.useContext(UserContext)

    return (
        <div style={{ padding: 20 }}>
            <Grid container spacing={5} justify="center">
                <Grid item xs={6}  >
                    <Card className={classes.card} variant="outlined">
                        <CardContent>
                            <Typography gutterBottom variant="h5" component="h2">
                                {props.title}
                            </Typography>
                            <Typography variant="body2" component="p">
                                <div>
                                    <pre style={{ whiteSpace: "pre-wrap" }}>{JSON.stringify(user, null, 4)}
                                    </pre>
                                </div>
                            </Typography>
                        </CardContent>
                        <CardActions>
                            <Button size="small">Learn More</Button>
                        </CardActions>
                    </Card>
                </Grid>
            </Grid>
        </div>
    );
}
