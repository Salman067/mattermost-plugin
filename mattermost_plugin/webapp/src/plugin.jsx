import React from 'react';

import {FormattedMessage} from 'react-intl';

import en from 'i18n/en.json';

import es from 'i18n/es.json';

import manifest from './manifest';

import RHSView from './components/right_hand_sidebar';

import {
    ChannelHeaderButtonIcon,
} from './components/icons';
import {
    channelHeaderMenuAction,
    websocketStatusChange,
    getStatus,
} from './actions';
import reducer from './reducer';

function getTranslations(locale) {
    switch (locale) {
    case 'en':
        return en;
    case 'es':
        return es;
    }
    return {};
}

export default class WFHPlugin {
    initialize(registry, store) {
        const {toggleRHSPlugin} = registry.registerRightHandSidebarComponent(
            RHSView,
            <FormattedMessage
                id='plugin.name'
                defaultMessage='WFH Plugin'
            />);

        registry.registerChannelHeaderButtonAction(
            <ChannelHeaderButtonIcon/>,
            () => store.dispatch(toggleRHSPlugin),
            <FormattedMessage
                id='plugin.name'
                defaultMessage='WFH Plugin'
            />,
        );

        registry.registerChannelHeaderMenuAction(
            <FormattedMessage
                id='plugin.name'
                defaultMessage='WFH Plugin'
            />,
            (channelId) => store.dispatch(channelHeaderMenuAction(channelId)),
        );

        registry.registerWebSocketEventHandler(
            'custom_' + manifest.id + '_status_change',
            (message) => {
                store.dispatch(websocketStatusChange(message));
            },
        );

        registry.registerReducer(reducer);

        store.dispatch(getStatus());

        registry.registerReconnectHandler(() => {
            store.dispatch(getStatus());
        });

        registry.registerTranslations(getTranslations);
    }

    uninitialize() {
        //eslint-disable-next-line no-console
        console.log(manifest.id + '::uninitialize()');
    }
}
