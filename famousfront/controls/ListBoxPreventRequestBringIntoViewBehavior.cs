using System.Windows;
using System.Windows.Controls;
using System.Windows.Interactivity;
using famousfront.utils;

namespace famousfront.controls
{
  class ListBoxPreventRequestBringIntoViewBehavior : Behavior<ListBox>
  {
    protected override void OnAttached()
    {
      base.OnAttached();
      TryHook();
      AssociatedObject.Loaded += new RoutedEventHandler(OnAssociatedObject).MakeWeakSpecial(eh=>AssociatedObject.Loaded -= eh);
    }
    bool _hooked;
    void OnAssociatedObject(object sender, RoutedEventArgs e)
    {
      if (_hooked)
        return;
      TryHook();
    }
    void TryHook()
    {
      var scp = VisualTreeExtensions.FindVisualChild<ScrollContentPresenter>(AssociatedObject);
      if (scp == null)
        return;
      _hooked = true;
      scp.RequestBringIntoView += new RequestBringIntoViewEventHandler(PreventRequestBringIntoView).MakeWeakSpecial(eh =>
      {
        scp.RequestBringIntoView -= eh;
      });
    }
    void PreventRequestBringIntoView(object sender, RequestBringIntoViewEventArgs e)
    {
      e.Handled = true;
    }
  }
}