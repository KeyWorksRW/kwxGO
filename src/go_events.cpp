/*
 * go_events.cpp - Event routing for kwxGO
 *
 * Provides wxEvtHandler connection/disconnection using a Go-compatible event
 * routing class. This follows the same pattern as kwxFortran's kwxApp.cpp
 * HandleEvent and SampleApp::HandleEvent in kwxFFI.
 *
 * Uses a callback registration pattern to avoid circular dependencies:
 * Go registers its callback function pointer via kwxGO_SetEventCallback()
 * before making any event connections.
 *
 * Flow:
 *   1. Go calls kwxGO_SetEventCallback(goEventCallback) at init
 *   2. Go calls kwxGO_Connect(obj, first, last, type, data)
 *   3. C++ creates wxClosure(g_goEventCallback, data) and wxCallback
 *   4. wxEvtHandler::Connect stores the callback
 *   5. Event fires → GoEventRouter::HandleEvent → wxCallback::Invoke
 *   6. wxClosure::Invoke calls g_goEventCallback(fun, data, evt)
 *   7. goEventCallback (Go) dispatches to the registered Go handler
 */

#include "kwx_wrapper.h"

#include <wx/wx.h>

// Function pointer type for Go event callback
typedef void (*GoEventCallbackFunc)(void *, void *, void *);

// Global function pointer - set by Go at runtime via kwxGO_SetEventCallback
static GoEventCallbackFunc g_goEventCallback = nullptr;

/*-----------------------------------------------------------------------------
    GoEventRouter - Event handler class for routing events to Go callbacks
-----------------------------------------------------------------------------*/

/// GoEventRouter provides the HandleEvent method pointer needed by
/// wxEvtHandler::Connect. This is the Go-specific equivalent of
/// SampleApp::HandleEvent from kwxFFI.
class GoEventRouter : public wxEvtHandler
{
public:
    void HandleEvent(wxEvent &evt)
    {
        wxCallback *callback = static_cast<wxCallback *>(evt.m_callbackUserData);
        if (callback && g_goEventCallback)
        {
            callback->Invoke(&evt);
        }
    }
};

/*-----------------------------------------------------------------------------
    C Interface for Go event connection
-----------------------------------------------------------------------------*/
extern "C"
{
    /// Register the Go event callback function pointer.
    /// Must be called before any kwxGO_Connect calls.
    ///
    /// @param callback Function pointer to Go's goEventCallback function
    EXPORT void kwxGO_SetEventCallback(GoEventCallbackFunc callback)
    {
        g_goEventCallback = callback;
    }

    /// Connect a Go callback to a wxWidgets event.
    ///
    /// @param obj   The wxEvtHandler to connect to (window, button, etc.)
    /// @param first First ID in range (-1 for any ID)
    /// @param last  Last ID in range (-1 for any ID)
    /// @param type  Event type (e.g., EVT_COMMAND_BUTTON_CLICKED)
    /// @param data  Opaque Go callback registry ID (uintptr cast to void*)
    /// @return 0 on success, -1 if callback not registered
    EXPORT int kwxGO_Connect(void *obj, int first, int last, int type, void *data)
    {
        if (!g_goEventCallback)
        {
            return -1; // Callback not registered yet
        }

        auto *closure = new wxClosure(g_goEventCallback, data);
        auto *callback = new wxCallback(closure);
        static_cast<wxEvtHandler *>(obj)->Connect(
            first, last, type,
            (wxObjectEventFunction)&GoEventRouter::HandleEvent, callback);
        return 0;
    }

    /// Disconnect all Go callbacks for a specific event type.
    ///
    /// @param obj   The wxEvtHandler to disconnect from
    /// @param first First ID in range (-1 for any ID)
    /// @param last  Last ID in range (-1 for any ID)
    /// @param type  Event type to disconnect
    /// @return 1 if a handler was disconnected, 0 if not found
    EXPORT int kwxGO_Disconnect(void *obj, int first, int last, int type)
    {
        return static_cast<int>(static_cast<wxEvtHandler *>(obj)->Disconnect(
            first, last, type,
            (wxObjectEventFunction)&GoEventRouter::HandleEvent));
    }
}
